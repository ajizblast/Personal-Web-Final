package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// kenapa struct tidak di deklarasi di dalam function main tapi di luar?
// agar nanti struct yg kita deklarasikan bisa dipakai di function lain

type Blog struct {
	ID              int
	Title           string
	StartDate       time.Time
	EndDate         time.Time
	Duration        string
	Content         string
	Author          string
	NodeJs          bool
	ReactJs         bool
	NextJs          bool
	TypeScript      bool
	Image           string
	FormatStartDate string
	FormatEndDate   string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var userData = SessionData{}

// membuat tempat peampung slice, yang berguna sebagai data dummy yg akan muncul di web
// kenapa kita menggunakan slice? karena di slice bisa menggunakan pengoperasian di slice

func main() {

	connection.DatabaseConnection()
	//create echo new instance
	e := echo.New()

	// serve static files from public directory(untuk mengakses direktory file)
	e.Static("/public", "public")
	e.Static("/uploads", "uploads")

	// To use sessions using eccho
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	//Routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/blog", blog)
	e.GET("/blog-detail/:id", blogDetail)
	e.GET("/edit-blog/:id", editBlog)
	e.GET("/testimonial", testimonial)

	e.POST("/add-blog", middleware.UploadFile(addBlog))
	e.POST("/edit-blog/:id", middleware.UploadFile(submitEditBlog))
	e.POST("/delete-blog/:id", deleteBlog)

	//register
	e.GET("/form-register", formRegister)
	e.POST("/register", register)

	//login
	e.GET("form-login", formLogin)
	e.POST("/login", login)

	e.POST("/logout", logout)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {

	// data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, content, nodejs, nextjs, reactjs, typescript, image FROM tb_blog")
	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_blog.id, title, start_date, end_date, content, nodejs, nextjs, reactjs, typescript, image, tb_user.name AS author FROM tb_blog JOIN tb_user ON tb_blog.author_id = tb_user.id ORDER BY tb_blog.id DESC")

	var result []Blog
	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.StartDate, &each.EndDate, &each.Content, &each.NodeJs, &each.NextJs, &each.ReactJs, &each.TypeScript, &each.Image, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}

		// each.Author = "Chahyo"
		each.FormatStartDate = each.StartDate.Format("2 January 2006")
		each.FormatEndDate = each.EndDate.Format("2 January 2006")

		result = append(result, each)
	}

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	datas := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"Blogs":        result,
		"DataSession":  userData,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func contact(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	datas := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"DataSession":  userData,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/contact-me.html")
	if err != nil {
		// fmt.Println("Tidak ada datanya")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func blog(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		// fmt.Println("Tidak ada datanya")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func blogDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //123 string => 123 int

	var blogDetail = Blog{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT tb_blog.id, title, start_date, end_date, content, nodejs, nextjs, reactjs, typescript, image, tb_user.name AS author FROM tb_blog JOIN tb_user ON tb_blog.author_id = tb_user.id WHERE tb_blog.id=$1", id).Scan(
		&blogDetail.ID, &blogDetail.Title, &blogDetail.StartDate, &blogDetail.EndDate, &blogDetail.Content, &blogDetail.NodeJs, &blogDetail.NextJs, &blogDetail.ReactJs, &blogDetail.TypeScript, &blogDetail.Image, &blogDetail.Author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// blogDetail.Author = "Chahyo Purnomo Aji"
	blogDetail.Duration = "0 Month"
	blogDetail.FormatStartDate = blogDetail.StartDate.Format("2 January 2006")
	blogDetail.FormatEndDate = blogDetail.EndDate.Format("2 January 2006")

	data := map[string]interface{}{
		"Blog":      blogDetail,
		"StartDate": getDateString(blogDetail.StartDate.Format("2006-01-02")),
		"EndDate":   getDateString(blogDetail.EndDate.Format("2006-01-02")),
	}

	var tmpl, errTemplate = template.ParseFiles("views/blog-detail.html")

	if errTemplate != nil {
		// fmt.Println("Tidak ada datanya")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func editBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var blogDetail = Blog{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, content, nodejs, nextjs, reactjs, typescript, image, FROM tb_blog WHERE tb_blog.id=$1", id).Scan(
		&blogDetail.ID, &blogDetail.Title, &blogDetail.StartDate, &blogDetail.EndDate, &blogDetail.Content, &blogDetail.NodeJs, &blogDetail.NextJs, &blogDetail.ReactJs, &blogDetail.TypeScript, &blogDetail.Image)

	data := map[string]interface{}{
		"Blog": blogDetail,
	}

	var tmpl, errTemplate = template.ParseFiles("views/edit-blog.html")

	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func testimonial(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimonial.html")
	if err != nil {
		// fmt.Println("Tidak ada datanya")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func addBlog(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	nodeJs := (c.FormValue("nodeJs") == "nodeJs")
	reactJs := (c.FormValue("reactJs") == "reactJs")
	nextJs := (c.FormValue("nextJs") == "nextJs")
	typescript := (c.FormValue("typescript") == "typescript")
	image := c.Get("dataFile").(string)
	sess, _ := session.Get("session", c)

	author := sess.Values["id"].(int)

	fmt.Println(title, content, nodeJs, reactJs, nextJs, typescript, image, author)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_blog (title, content, start_date, end_date, nodejs, reactjs, nextjs, typescript, image, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", title, content, startDate, endDate, nodeJs, reactJs, nextJs, typescript, image, author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func submitEditBlog(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	nodeJs := (c.FormValue("nodeJs") == "nodeJs")
	reactJs := (c.FormValue("reactJs") == "reactJs")
	nextJs := (c.FormValue("nextJs") == "nextJs")
	typescript := (c.FormValue("typescript") == "typescript")
	image := c.Get("dataFile").(string)
	sess, _ := session.Get("session", c)

	author := sess.Values["id"].(int)

	fmt.Println(title, content, nodeJs, reactJs, nextJs, typescript, image, author)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE INTO tb_blog SET title=$1, content=$2, start_date=$3, end_date=$4, nodejs=$5, reactjs=$6, nextjs=$7, typescript=$8, image=$9, author_id=$10, WHERE id=$11", title, content, startDate, endDate, nodeJs, reactJs, nextJs, typescript, image, author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("ID: ", id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_blog WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Register
func formRegister(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error {
	// to make sure request body is form data format, not JSON, XML, etc.
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("inputName")
	email := c.FormValue("inputEmail")
	password := c.FormValue("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		redirectWithMessage(c, "Register failed, please try again.", false, "/form-register")
	}

	return redirectWithMessage(c, "Register success!", true, "/form-login")
}

// Login
func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	email := c.FormValue("inputEmail")
	password := c.FormValue("inputPassword")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "Email Incorrect!", false, "/form-login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Incorrect!", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 JAM
	sess.Values["message"] = "Login success!"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func getDateString(date string) string {
	y := date[0:4]
	m, _ := strconv.Atoi(date[5:7])
	d := date[8:10]
	if string(d[0]) == "0" {
		d = string(d[1])
	}

	mon := ""
	switch m {
	case 1:
		mon = "Jan"
	case 2:
		mon = "Feb"
	case 3:
		mon = "Mar"
	case 4:
		mon = "Apr"
	case 5:
		mon = "Mei"
	case 6:
		mon = "Jun"
	case 7:
		mon = "Jul"
	case 8:
		mon = "Agu"
	case 9:
		mon = "Sep"
	case 10:
		mon = "Okt"
	case 11:
		mon = "Nov"
	case 12:
		mon = "Des"
	}

	return d + " " + mon + " " + y
}
