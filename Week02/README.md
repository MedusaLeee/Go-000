## 作业

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 思路


1. `sql.ErrNoRows` 只是找不到对应的记录，在dao层不算sql错误
2. 应该在dao层捕获这个错误，返回nil，其他dao层的错误，wrap后返回
3. service层根据实际情况判断，是否报错或打印日志，或者其他处理，见下方举例

## 举例

1. 比如put /users/:id，id对应的user找不到就要报400错误了
2. 比如get /users/:id/exists，这样的话就不需要报错了，返回 true or false

## 伪代码
    // error
    NoRowsError := errors.New("sql.ErrNoRows")
    
    // model
    type User struct {}
    
    // dao 
    func FindUserById(userId string) (*model.User, error) {
        user := &model.User{}
        if err := db.Where(`"id" = ?`, userId).Find(user).Error; err != nil {
            if err == NoRowsError {
                return nil, nil
            }
            return nil, errors.Warp(err, fmt.Sprintf("dao.FindUserById sql error, userId: %s", userId))
        }
        return user, nil
    }
    
    // service
    func UserExists(userId string) (bool, error) {
        user, err := dao.FindUserById(userId)
        if err != nil {
            return false, error
        }
        if user != nil {
            return true, nil
        }
        return false, nil
    }
    
    // controller
    func UserExists(c *gin.Context) {
        exists, err := service.UserExists(userId)
        if err != nil {
            log.Error(err)
            ctx.JSON(500, gin.H{
                "errNo": "500",
                "msg":   "系统繁忙",
            })
            return
        }
        ctx.JSON(200, gin.H{
            "exists": exists
        })
    }