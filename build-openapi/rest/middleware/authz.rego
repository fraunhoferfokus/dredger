package simpleApi.authz

#
# Methods
#

default getMethod = false

getMethod {
	lower(input.method) == "get"
}

default putMethod = false

putMethod {
	lower(input.method) == "put"
}

default postMethod = false

postMethod {
	lower(input.method) == "post"
}

default deleteMethod = false

deleteMethod {
	lower(input.method) == "delete"
}

#
# Roles
#

default staff = false

staff {
	lower(input.role) == "staff"
}

default user = false

user {
	lower(input.role) == "user"
}

default staffuser = false

staffuser {
	staff
}

staffuser {
	user
}

#
# API token
#

default api = false

api {
	lower(input.role) == "api"
}

api {
	input.apitoken == ""
}

#
# Rules
#
# allowEntryPoint permit general the access to a api function
# allowAccess check the authorization by the permitted roles
#
default allowEntrypoint = false

default allowAccess = false

#
# examples how you could set permission for each path in the API endpoint
#

allowEntrypoint {
    getMethod
    input.path == "/livez"
}

allowAccess {
	getMethod
    input.path == "/livez"
}

allowEntrypoint {
    getMethod
    input.path == "/readyz"
}

allowAccess {
	getMethod
    input.path == "/readyz"
}

allowEntrypoint {
    getMethod
    input.path == "/infoz"
}

allowAccess {
	getMethod
    input.path == "/infoz"
}

allowEntrypoint {
    getMethod
    input.path == "/robots.txt"
}

allowAccess {
	getMethod
    input.path == "/robots.txt"
}

allowEntrypoint {
    user
#	getMethod
#   input.path == "/func"
}

allowAccess {
	user
#	getMethod
#   input.path == "/func"
}

#allowEntrypoint {
#   staff
#	postMethod
#   input.path == "/admin"
#}

#allowAccess {
#   staff
#	postMethod
#   input.path == "/admin"
#}

#allowEntrypoint {
#   api
#	deleteMethod
#   input.path == "/func"
#}

#allowAccess {
#   api
#	deleteMethod
#   input.path == "/func"
#}
