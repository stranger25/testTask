create table R_CURRENCY
(
	ID int auto_increment,
	TITLE varchar(60) null,
	CODE varchar(3) null,
	VALUE numeric(18,2) null,
	A_DATE date null,
	constraint R_CURRENCY_pk
		primary key (ID)
);

