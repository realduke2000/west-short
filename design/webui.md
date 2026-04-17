# Short Server Management Web UI

## Basics:
1. Use react framework to build webui

## Pages

### Login Page

#### Overview
This page accept username and password input, authenticate user identity and permissions, and redirect user to corresponding page.

#### Layout
URL: /login 
Title: Short Server Management Console

Input:
* User name text box - [a-z][A-Z][._-][0-9] are allowed
* password textbox
  * mask user input 
* Login Button

#### Events:
Behavior:  
1. If user input correct username and password, redirect to management page, url "/manage"
2. if user input incorrect username or password, display error message, note: 
   1. Do not prompt whether username or password is incorrect

User authentication workflow:
1. Accept user name and password
3. Make sure user name only contains legal chars
4. call backend authentication API, request body is a json string `{"username":"xxx", "password":"xxx"}`, use sha256 to hash plain text password, and then use base64 encoding to encode plaintext password, format is: "sha256:<base64 encoded sha256 password>"
1. If API returns 200 OK, read token from response body, and then store in http-only local cookie, expires in 24 hours.
1. If API returns non-ok code, display login error message on webui.

### Admin page

### Management page

#### Overview
User can manage short url here, all short urls created by current user are listed in this page, user can create new short url, delete short url in this page.

Display scenario:  
1. Use a paginated grid to display all short urls
1. Column head should be short url, which is id, and long url, last access time, created time.

Logout scenario:  
1. On top-right corner, displays user icon, when mouse over the icon, pop up a drop list, logout button is over there.
1. When click logout button, clear token cookie and redirect page to default page "/"
