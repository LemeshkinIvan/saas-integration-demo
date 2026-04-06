package source_jobs

import (
	"context"
	inst "daos_core/internal/data/repositories/instance"
	dto "daos_core/internal/domain/dto/instance"
	source_adapter "daos_core/internal/external/adapters/source"
	custom_logger "daos_core/internal/utils/logger"
	"fmt"
	"time"
)

type UpdateSourceJob struct {
	InstanceID int64
	Data       dto.UpdateDTO
	Access     string
	Referer    string
	Repo       inst.Repository
	Adap       source_adapter.Adapter
}

func (j *UpdateSourceJob) Do() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := j.Adap.UpdateByID(j.Referer, j.Access, j.Data)
	if err != nil {
		custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] AmoCRM update source failed", j.Data.InstanceID))
		return j.Repo.SetAmoError(ctx, j.Data.InstanceID, err.Error())
	}

	// if err := j.Repo.Update(
	// 	ctx,
	// 	instance.UpdateInput{
	// 		InstanceID: j.Data.InstanceID,
	// 		AccountID:  j.Data.AccountID,
	// 		Name:       j.Data.Name,
	// 		PipelineID: j.Data.PipelineID,
	// 		SourceID:   j.Data.SourceID,
	// 	}); err != nil {
	// 	custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] DB update source failed", j.Data.InstanceID))
	// 	return j.Repo.SetAmoError(ctx, j.Data.InstanceID, err.Error())
	// }

	if err := j.Repo.DeleteAmoErr(ctx, j.Data.InstanceID); err != nil {
		custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] %s", j.Data.InstanceID, err.Error()))
		return err
	}

	//custom_logger.AsyncLog(1, fmt.Sprintf("[Instance %d] succsessfully updated", j.Data.InstanceId))
	return nil
}
