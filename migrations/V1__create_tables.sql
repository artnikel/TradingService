create table deal (
	dealid uuid,
	profileid uuid,
	company varchar,
	purchaseprice double precision,
	sharescount double precision,
	stoploss double precision,
	takeprofit double precision,
	dealtime timestamp,
	enddealtime timestamp,
	profit double precision,
	primary key (dealid)
);