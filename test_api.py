''' 一个简单的测试程序 '''
import requests
import json
import logging
import pytest
import sqlite3
import random
import string

@pytest.fixture
def test_readConfig(filename: str='config.json') -> dict:
    ''' 读取配置文件 返回字典 或者空字典'''
    try:
        with open(filename, 'r') as f:
            config = json.load(f)
    except FileNotFoundError:
        logging.error('配置文件不存在')
        return {}
    except json.decoder.JSONDecodeError:
        logging.error('配置文件格式错误')
        return {}

    # 检查必要参数 是否存在
    if 'db' not in config:
        logging.error('配置文件缺少url参数')
        return {}
    elif 'host' not in config:
        logging.error('配置文件缺少host参数')
        return {}
    elif 'port' not in config:
        logging.error('配置文件缺少port参数')
        return {}
    return config

@pytest.fixture
def test_get_user(test_readConfig):
    ''' 随机获得一个用户 '''
    conn = sqlite3.connect(test_readConfig['db'])
    cursor = conn.cursor()
    cursor.execute('SELECT username, password FROM users ORDER BY RANDOM() LIMIT 1')
    user = cursor.fetchone()
    cursor.close()
    conn.close()
    return user

def test_ping(test_readConfig):
    ''' 测试ping接口 '''
    url = 'http://{host}:{port}/ping'.format(**test_readConfig)
    resp = requests.get(url, timeout=1)
    assert resp.status_code == requests.codes.ok
    respJson = resp.json()
    assert respJson['code'] == 0
    assert respJson['msg'] == 'pong!'

class TestSignin:
    ''' 测试登录接口 '''

    @pytest.fixture(autouse=True)
    def setup(self, test_readConfig):
        self.url = 'http://{host}:{port}/signin'.format(**test_readConfig)

    def test_signin(self, test_readConfig, test_get_user):
        ''' 测试正常登录 '''
        user = test_get_user
        data = {'username': user[0], 'password': user[1]}

        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] == 0
        assert respJson['msg'] == ''
    
    def test_signin_wrong_password(self, test_readConfig, test_get_user):
        ''' 测试输入错误的密码 '''
        user = test_get_user
        data = {'username': user[0], 'password': 'kk'}

        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] != 0

    def test_signin_no_username(self, test_readConfig):
        ''' 测试没有用户名登录 '''
        data = {'password': '123456'}
        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] != 0

class TestCheckin:
    ''' 测试签到接口 '''
    @pytest.fixture(autouse=True)
    def setup(self, test_readConfig, test_get_user):
        self.url = 'http://{host}:{port}/checkin'.format(**test_readConfig)
        # Get token
        user = test_get_user
        data = {'username': user[0], 'password': user[1]}
        resp = requests.post('http://{host}:{port}/signin'.format(**test_readConfig), json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] == 0
        self.token = respJson['data']['access_token']
    
    def test_checkin(self):
        ''' 测试签到 '''
        data = {
            'access_token': self.token,
            'word': 'test',
            }
        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] == 0
        assert type(respJson['data']['point']) == int
    
    def test_checkin_no_token(self):
        ''' 测试没有token '''
        data = {'word': 'test'}
        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] != 0

    def test_checkin_invalid_token(self):
        ''' 测试无效的token '''
        data = {
            'access_token': 'invalid',
            'word': 'test',
            }
        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] != 0
    
class TestSignup:
    ''' 测试注册接口 '''
    @pytest.fixture(autouse=True)
    def setup(self, test_readConfig, test_get_user):
        self.url = 'http://{host}:{port}/signup'.format(**test_readConfig)
        # get a known user
        self.user = test_get_user
    
    def test_signup(self):
        ''' 测试注册 '''
        # 随机的8位用户名和密码[0-9 a-z A-Z]
        username = ''.join(random.choices(string.ascii_letters + string.digits, k=8))  
        password = ''.join(random.choices(string.ascii_letters + string.digits, k=8))
        data = {
            'username': username,
            'password': password,
        }
        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] == 0

    def test_signup_with_exist_username(self):
        ''' 测试注册已经存在的用户名 '''
        data = {
            'username': self.user[0],
            'password': self.user[1],
        }
        resp = requests.post(self.url, json=data, timeout=1)
        assert resp.status_code == requests.codes.ok
        respJson = resp.json()
        assert respJson['code'] != 0
