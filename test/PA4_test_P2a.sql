USE CS457_PA4;
select * from Flights;
begin transaction;
update flights set status = 1 where seat = 22;
commit; --there should be nothing to commit; it's an "abort"
select * from Flights;
