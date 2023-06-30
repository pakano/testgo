#!/usr/local/bin/lua

-- HTTP header
print("Content-Type: text/html\n\n")
print("")          -- An empty line

-- body
print("<h1>Hello CGI</h1>")
print("<h2>Current time: " .. os.date("%D %T") .. "</h2")

