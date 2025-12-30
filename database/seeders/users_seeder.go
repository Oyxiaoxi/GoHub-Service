package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/console"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/seed"

	"gorm.io/gorm"
)

func init() {

	// 添加 Seeder
	seed.Add("SeedUsersTable", func(db *gorm.DB) {

		rand.Seed(time.Now().UnixNano())

		cities := []string{"Shanghai", "Beijing", "Shenzhen", "Guangzhou", "Hangzhou", "Chengdu", "Wuhan", "Nanjing", "Xiamen", "Remote"}
		intros := []string{
			"爱折腾后端与基础设施，周末骑行", "关注 Go 微服务与观测性", "业余摄影师，喜欢人文题材",
			"前端转后端，正在学习分布式", "移动端开发，偶尔写服务端", "安全向，热爱代码审计与红蓝对抗",
			"产品经理转码农，记录踩坑", "全栈打杂，关注性能与用户体验", "云原生爱好者，K8s 练习生",
			"学生党，刷题与开源新手",
		}

		var users []user.User
		highActive := 5
		moderate := 10
		low := 15

		addUser := func(idx int, activeLevel string) {
			city := cities[rand.Intn(len(cities))]
			intro := intros[rand.Intn(len(intros))]
			avatar := ""
			// 一部分用户保留空头像，其余使用示例头像链接
			if rand.Intn(100) < 65 {
				avatar = fmt.Sprintf("https://example.com/avatars/%02d.png", idx)
			}

			u := user.User{
				Name:           fmt.Sprintf("user_%s_%02d", activeLevel, idx),
				Email:          fmt.Sprintf("user_%s_%02d@example.com", activeLevel, idx),
				Phone:          fmt.Sprintf("1%010d", rand.Intn(9_000_000_000)+1_000_000_000),
				Password:       "$2a$14$oPzVkIdwJ8KqY0erYAYQxOuAAlbI/sFIsH0C0R4MPc.3JbWWSuaUe", // 123456
				City:           city,
				Introduction:   intro,
				FollowersCount: int64(rand.Intn(200)),
				Points:         int64(rand.Intn(500)),
				Avatar:         avatar,
			}
			users = append(users, u)
		}

		for i := 0; i < highActive; i++ {
			addUser(i+1, "hi")
		}
		for i := 0; i < moderate; i++ {
			addUser(i+1, "mid")
		}
		for i := 0; i < low; i++ {
			addUser(i+1, "lo")
		}

		result := db.Table("users").Create(&users)
		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.RowsAffected))
	})
}
