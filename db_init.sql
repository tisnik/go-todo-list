create table todo (
    ID            integer primary key asc,
    due_to        text not null,
    finished_at   text,
    priority      number not null,
    subject       text not null,
    details       text not null
);

insert into todo (id, due_to, finished_at, priority, subject, details) values (0, '2024-12-31', null, 1, 'first item', 'some details');
insert into todo (id, due_to, finished_at, priority, subject, details) values (1, '2024-12-31', '2024-10-25', 5, 'start Go workshop', 'problems with sharing the screen (as usual)');
