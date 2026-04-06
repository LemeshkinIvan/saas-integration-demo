package sqlmodels

import "time"

// create table instances(
// 	id serial not null primary key,
// 	account_id int,
// 	bot_token_id int,
// 	name varchar(150) not null default 'undefined',
// 	created_at timestamp default NOW(),
// 	status int not null default 2,
// 	updated_at timestamp,

// 	FOREIGN KEY (bot_token_id) REFERENCES bots_token (id),
// 	FOREIGN KEY (account_id) REFERENCES accounts (id)
// );

// пока хз мб для тестов
type Instances struct {
	Id int64 `db:"id"`
	// подразумеваю что здесь будет уже значения после join, а не id таблиц, где serial
	IdUserTokens string    `db:"id_user_tokens"`
	IdBotToken   string    `db:"id_bot_token"`
	Name         string    `db:"name"`
	CreatedAt    time.Time `dn:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
