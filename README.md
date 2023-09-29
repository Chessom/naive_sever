# naive_sever
A simple sever written with Golang
## Introduction
This program use Gin,the web framework, to deal with the work with web, http and json parsing.It also use GORM to operate the database.
This program implements the following functions: registration, login, and sign-in. The identity authentication token uses a randomly 
generated uuid. Signing in can record points in the database, and you can only sign in once a day.

## The detailed error codes' meanings
Code|meaning
0|Done
1|Empty data
2|not existent
3|incorrect password
4|invalid access token
5|repeated checkin
6|username has existed
7|internal unknown error

问题回答：
1.有意思（虽然Gin和GORM的一些大小写规则真的很奇怪）
2.文档里给出的资料就基本足够，我自己找了go语言的官方文档和一些关于gin、GORM的博客
3.困难就是不太适应新的语言语法习惯（但我觉得多用用可能就习惯了），以及一些古怪的字段命名规则让我不知错误在何处。
前者就是时间会让人适应一切，后者则是search the web解决的疑难杂症。
4.没有了。我觉得可以。
5.无
