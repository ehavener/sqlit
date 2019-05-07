CREATE DATABASE CS457_PA4;
USE CS457_PA4;
create table Flights(seat int, status int);

insert into Flights values(22,0);

select * from Flights;

begin transaction;

update flights set status = 1 where seat = 22;

commit;

select * from Flights;

.EXIT