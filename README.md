# onthemat_backend

# 1. 개요

✔️ 서비스 : 요가 대강 매칭 서비스

계속 진행중 .....



## 1.1. 스택

🔎 백엔드  
<div style="display:flex;">
   <img src="https://img.shields.io/badge/GO-gray?style=flat&logo=Go&logoColor=00ADD8"/>
	<img src="https://img.shields.io/badge/Fiber-white?style=flat"/>
	<img src="https://img.shields.io/badge/EntGO-white?style=flat"/>
</div>
<br>

🔎 데이터베이스  
<div style="display:flex;">
  <img src="https://img.shields.io/badge/PostgreSQL-green?style=flat&logo=PostgreSQL&logoColor=4169E1"/>
  <img src="https://img.shields.io/badge/Redis-green?style=flat&logo=Redis&logoColor=DC382D"/>
</div>
<br>


## 1.2. 디렉토리 구조

```
├── cmd
│   ├── app             // 메인 어플리케이션
│ . ├── migraiton				// 마이그레이션 생성
│   └── seeding 				// 테스트용 데이터 insert
│
│
├── configs             // 어플리케이션 설정파일
│
├── internal
│    └── app 
│         ├── delivery            
│         │           ├── http   // http handler
│         │           └── middleware 
│         ├── config
│         │
│         ├── common
│         │
│         ├── infrastructure
│         │
│         ├── migrations
│         │
│         ├── mocks
│         │
│         ├── model
│         │
│         ├── repository
│         │
│         ├── service
│         │
│         ├── usecase
│         │
│         ├── transport // 데이터 전송 Object
│         │           ├── request
│         │           └── response
│         │
│         ├── utils
│
├── pkg
│    ├── auth 
│    │      ├── jwt
│    │      └── store
│    │
│    ├── aws
│    │
│    ├── email
│    │
│    ├── ent
│    │
│    ├── entx
│    │
│    ├── fiberx
│    │
│    ├── google
│    │
│    ├── mocks
│    │
│    ├── kakao
│    │
│    ├── naver
│    │
│    ├── openapi
│    │
│    ├── validatorx
│
│
├── scripts
```

## 1.3. 아키텍처


### 1.3.1. 서버 아키텍처

의존성을 최대한 낮추기 위해 노력했습니다.

![서버 아키텍처](https://user-images.githubusercontent.com/97140962/202600851-884abaad-c12c-4f7e-8b23-715dee475e5c.jpg)

## 1.4. API 명세 

URL : http://43.201.147.22:3000/

![api문서](https://user-images.githubusercontent.com/97140962/201019708-08588b56-8304-4a77-946a-cf67e443a7a5.png)



# 2. 기술 상세

## 2.1. REST API

RESTFUl한 API를 설계했습니다.

GET, POST, PUT, PACTH, DELETE 5가지를 사용합니다.

### PUT
db에 저장된 값을 요청 값으로 모두 대체합니다.
만약 저장된 값이 없을 경우 요청 값을 생성합니다.

### Patch

Put메서드와 원리는 같으나

Put은 저장된 값 전부를 대체하지만,
Patch는 요청한 값만 대체합니다.

만약 저장된 값이 없을 경우 요청값을 생성하며,

데이터베이스에 필수로 들어가야 하는 값이
전부 있는 경우에만 성공적인 응답을 받을 수 있습니다.


## 2.2. Repository

### 2.2.1. 업데이트 로직 

[업데이트 로직 예시 ](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L112)

데이터베이스에 컬럼이 NULL이 가능한 경우는
요청값이 없으면 NULL로 변경합니다.

다대다 관계에서는
기존 값을 전부 지우고 새로운 값으로 대체합니다.


[일대다 관계일 때 로직](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L153)
일대다 관계에서는 
데이터베이스에 존재하는 id값들과 요청값 id들을 비교하여
생성, 업데이트, 삭제를 진행합니다.


### 2.2.2. Patch 로직

[Patch로직 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L292)


요청값에 id가 있으면 업데이트 없으면 생성,

요청값이 업데이트 가능한 컬럼이면
업데이트 합니다.

### 2.2.3. List 로직

[List로직 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/recruitment.go#L291)

위 부분은 PostgreSQL의 JsonB타입 컬럼에 존재하는
데이터 Key값을 기준으로 리스트를 조회할 수 있도록 도와주는
검색로직입니다.


### 2.2.4. 테스트케이스


[Repository_test](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/user_test.go)

테스트 시작 시, 로컬에서 도커 컨테이너를 생성하고 
테스트 전체가 종료되면
컨테이너를 삭제하도록 하였습니다.

각각의 테스트케이스마다 독립적인 실행을 원했습니다.
각각의 테스트가 종료되면 테스트한 데이터 전체를 삭제하여
의존성을 낮췄습니다.

각각의 테스트케이스마다  
테스트를 위한 데이터를 미리 생성해야되기 때문에,  
(ex, SELECT 을 했을 때 정상적으로 데이터가 출력됨을 확인하기 위해 INSERT작업이 필요)

테스트 코드가 길어져 가독성이 떨어진다는 단점이 존재했습니다.  

BeforeTest 함수에 테스트에 필요한 코드들을 넣음으로써
테스트케이스에는 테스트코드만 넣어 이해하기 쉽게 만들었습니다.


## 2.3. Service

이메일 전송, 엑셀 데이터 파싱 등 데이터베이스와 무관한 코드들의 모음입니다.
### 2.3.1. 행정구역 엑셀 파싱

[행정구역 엑셀 파싱](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/service/area_service.go#L31)

엑셀의 행정구역 데이터를 파싱하는 로직입니다.



## 2.4. Usecase

### 2.4.1. 인증

소셜로그인, 토큰 재발급, 회원가입 로그인, 이메일 인증, 비밀번호 초기화 등을 구현했습니다.

엑세스토큰은 만료시간을 15분으로 짧게하였고
Refresh는 2주로, Redis에 저장합니다.

사용자는 주기적으로 Refresh토큰을 가지고
Access토큰을 재발급하여 인증에 사용합니다.

[토큰을 재발급 받는 Usecase](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/usecase/auth_usecase.go#L372)

### 2.4.2. 생성

[Create](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/usecase/recruitment_usecase.go#L34)

Ent라는 ORM은 코드를 Generate해줍니다.

데이터베이스 스키마에 맞게
Object를 생성해줍니다. 

생성된 Object를 이용하여 사용자의 HTTP 요청을 받을 수도 있었지만 

데이터를 Request Response할 때의 Object와
데이터를 INSERT GET할 때의 Object를 분리하였습니다.

코드가 길어진다는 단점이 존재하지만,

요청값과 응답값의
코드를 작성함으로써
요청값과 응답값이 명확하다는 장점이 존재했고

요청 응답 값을 바꿔도 데이터베이스에 접근하는 오브젝트를 수정하지 않아도 되서 편리했습니다.

### 2.4.3. 에러처리

[에러 서버코드](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/common/errors.go#L50)

HTTP응답코드를 준수하였으며,
400, 401, 403, 400, 409 등이 사용됩니다.

친절한 API를 만들고 싶었습니다.
에러메세지와 에러코드를 상황별로 제공하여
프론트엔드 개발자가 예외처리를 조금 더 편하게 할 수 있도록 하였습니다.

### 2.4.4. 테스트케이스

[테스트케이스 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/usecase/user_usecase_test.go#L28)

Usecase에서는 repostiroy, service모듈이 사용됩니다.
사용되는 모듈들은 테스트케이스를 통해 검증이 완료된 모듈입니다.

따라서 usecase에서는 테스트가 불필요합니다.

mock을 사용하여 각 모듈의 리턴값을 원하는 값으로 정하여
불필요한 과정들을 거치지 않아도 되게끔 하였습니다.


## 2.4. handler

[코드 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/delivery/http/academy_handler.go#L30)

응답받은 요청을 파싱, 검증하는 과정을 담당합니다.