package dbrepo

import (
	"database/sql"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/repository"
)

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// Return new repo for postgres database
func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// Return new repo for postgres database
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

/* For add mySQL
Also go to handlers.go and fix NewRepo func
type mysqlDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewMySQLRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &mysqlDBRepo{
		App: a,
		DB:  conn,
	}
} */

/* some advantage of repository pattern
1. thứ nhất bạn tạo ra một list các method làm việc với logic của handler trong repo -> nhìn phát là kiểm soát luôn code mình sẽ làm những gì ?

2. repo_impl trong trường hợp này bạn implement cho postgres, nhưng nếu bạn muốn implement cho mysql thì sao. Cách giải quyết -> bạn vẫn sử dụng được repo interface và chỉ thêm 1 class mới ví dụ mysql_repo_impl và khi chuyển từ postgres -> mysql thì mình chỉ cần chỉnh lại thay vì dùng postgres_repo_impl thì mình dùng mysql_repo_impl, code chỉ đổi 1 chỗ.

3. dễ viết mock test hơn */
