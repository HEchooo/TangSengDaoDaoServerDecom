alter table user_info add column im_username varchar(255)  default '' not null comment 'IM登录用户名' after password;
alter table user_info add column im_password varchar(255)  default '' not  null comment 'IM登录密码'  after password;

