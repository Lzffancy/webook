# 接口文档

## 1.编辑用户信息

```
127.0.0.1:8082/users/edit
```

POST 

```
{
    "nickname": "fancy",
    "brithday": "2023-10-14",
    "aboutMe":"hhahah"
}
```

return

```
Edit is ok
```

![image-20231216002626021](C:\Users\fancy\AppData\Roaming\Typora\typora-user-images\image-20231216002626021.png)

## 2.显示用户信息

profile

```
127.0.0.1:8082/users/profile
```

GET
return

```
{
    "nickname": "fancy",
    "email": "test1@qq.com",
    "aboutMe": "hhahah",
    "birthday": "2023-10-14"
}
```

![image-20231217133941296](C:\Users\fancy\AppData\Roaming\Typora\typora-user-images\image-20231217133941296.png)