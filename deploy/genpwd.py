#!/usr/bin/env python3
# coding=utf-8
import hmac
import random
import string

rootpwd = ''.join(random.SystemRandom().choice(string.ascii_uppercase + string.ascii_lowercase + string.digits) for _ in range(16))
token = ''.join(random.SystemRandom().choice(string.ascii_uppercase + string.ascii_lowercase + string.digits) for _ in range(20))
pwdtoken = hmac.new(bytes(token, encoding = "utf-8"), bytes(rootpwd, encoding = "utf-8"), digestmod = 'SHA256').hexdigest()
print(pwdtoken, rootpwd, token)
