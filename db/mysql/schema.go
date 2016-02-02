package mysql

var (
	labelSchema = `CREATE TABLE IF NOT EXISTS labels (
id varchar(36) primary key,
service varchar(1024),
version varchar(255),
weight int(11),
priority int(11),
label varchar(255),
value varchar(255),
index(service),
index(priority),
index(label));`
)
