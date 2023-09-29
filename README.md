# naive_sever
A simple sever written with Golang
## Introduction
This program use Gin,the web framework, to deal with the work with web, http and json parsing.It also use GORM to operate the database.
This program implements the following functions: registration, login, and sign-in. The identity authentication token uses a randomly
generated uuid. Signing in can record points in the database, and you can only sign in once a day.

## The detailed error codes' meanings
|Code|meaning|
|------|------|
|0|Done|
1|Empty data|
|2|not existent|
|3|incorrect password|
|4|invalid access token|
|5|repeated checkin|
|6|username has existed|
|7|internal unknown error|

问题回答：  
1.有意思（虽然Gin和GORM的一些大小写规则真的很奇怪）  
2.文档里给出的资料就基本足够，我自己找了go语言的官方文档和一些关于gin、GORM的博客  
3.困难就是不太适应新的语言语法习惯（但我觉得多用用可能就习惯了），以及一些古怪的字段命名规则让我不知错误在何处。前者就是时间会让人适应一切，后者则是search the web解决的疑难杂症。  
4.没有了。我觉得可以。  
5.无  

感想:感觉回到了高中时那个激情澎湃的状态，不断地试错，不断地改进，不断地获得成就感。挫败与困惑在学习的过程中无法避免，但是在实践所学而成功的高光时刻才会永久留在心中。面对新鲜的语言，没用过的框架，完全不同的思考模式，我有着将它们一一收入囊中的渴望与热情。尽管时间有限，还是尽量完成了两题，这次经历很有意义，我诚恳地表达我对进入求是潮产品研发部的愿望，希望在这里可以遇到更多志同道合的朋友。
