package sqlmodels

// create table bot_token(
// 	id serial not null primary key,
// 	token text not null
// );

type BotToken struct {
	Id    int64  `db:"id"`
	Token string `db:"token"`
}
