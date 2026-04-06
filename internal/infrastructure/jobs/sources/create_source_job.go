package source_jobs

import (
	"context"
	"daos_core/internal/data/repositories/instance"
	source_adapter "daos_core/internal/external/adapters/source"
	custom_logger "daos_core/internal/utils/logger"
	"fmt"
	"time"
)

type CreateSourceJob struct {
	InstanceId int64
	Referer    string
	Access     string
	ExternalId string
	Repo       instance.Repository
	Adap       source_adapter.Adapter
}

func (j *CreateSourceJob) Do() error {
	//fmt.Println("check" + j.Referer + j.Access + j.ExternalId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	source, err := j.Adap.Create(
		j.Referer,
		j.Access,
		j.ExternalId,
	)
	if err != nil {
		custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] AmoCRM create source failed", j.InstanceId))
		return j.Repo.SetAmoError(ctx, j.InstanceId, err.Error())
	}

	if source == nil {
		custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] AmoCRM returned nil source", j.InstanceId))
		return j.Repo.SetAmoError(ctx, j.InstanceId, "source nil from AmoCRM")
	}

	if len(source.Embedded.Sources) == 0 {
		return fmt.Errorf("CreateSourceJob: Do: no sources returned")
	}

	//TODO
	updated := time.Now()
	err = j.Repo.UpdateBySource(ctx, j.InstanceId, source.Embedded.Sources[1], updated)
	if err != nil {
		custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] DB update after Amo failed: %v", j.InstanceId, err))
		return j.Repo.SetAmoError(ctx, j.InstanceId, err.Error())
	}

	if err := j.Repo.DeleteAmoErr(ctx, j.InstanceId); err != nil {
		custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] %s", j.InstanceId, err.Error()))
		return err
	}

	//custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] Source successfully created and updated in DB", j.InstanceId))
	return nil
}
