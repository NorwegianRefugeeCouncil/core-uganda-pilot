package main

//)
//
//func main() {
//	host := "localhost"
//	port := 5435
//	user := "postgres"
//	password := "postgres"
//	dbname := "core"
//	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
//	srv := bla.NewServer(psqlconn)
//	done := make(chan struct{})
//	if err := srv.Start(done, "127.0.0.1:9000"); err != nil {
//		panic(err)
//	}
//	<-done
//}
