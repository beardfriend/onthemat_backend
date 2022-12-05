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

fiber를 선택한 이유 : 압도적인 성능 때문입니다.

![image](https://user-images.githubusercontent.com/97140962/205558298-df3012cd-5f72-43a6-a158-1d987105198c.png)

entGO를 선택한 이유 : generate된 인터페이스를 사용하기 떄문에 gORM에 비해 안정적이면서도 빠른 속도를 가져갈 수 있습니다. 




## 1.2. 디렉토리 구조

```
├── cmd
│   ├── app             // 메인 어플리케이션
│   ├── migraiton	// 마이그레이션 생성
│   └── seeding 	// 테스트용 데이터 insert
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
│         ├── transport 	// 데이터 전송 Object
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



![서버 아키텍처](https://user-images.githubusercontent.com/97140962/202600851-884abaad-c12c-4f7e-8b23-715dee475e5c.jpg)

## 1.4. API 명세 

URL : http://43.201.147.22:3000/

<img width="1680" alt="스크린샷 2022-12-04 오후 3 54 02" src="https://user-images.githubusercontent.com/97140962/205478479-f0c7fd16-e8fd-4590-81f8-a1b8c386cfa6.png">


# 2. 기술 상세

## 2.1. REST API


GET, POST, PUT, PACTH, DELETE 5가지를 사용합니다.

RFC문서의 내용에 맞게 설계했습니다.

> The PUT method requests that the state of the target resource be created or 
> replaced with the state defined by the representation enclosed in the request message payload.

> This specification defines the new HTTP/1.1 [RFC2616] method, 
> PATCH, which is used to apply partial modifications to a resource.


PATCH는 사용자가 원하는 자원만 골라서    
수정할 수 있습니다.

PATCH를 한 번 만들어 놓는다면,  
클라이언트 측에서 부분적으로 수정할 요소들이 계속 변해도  
API를 수정하지 않아도 되는 장점이 있었습니다.  

이는 Graphql과 닮았습니다.  
Graphql은 조회할 때마저도 사용자가 원하는 자원을 가져올 수 있습니다.  
서비스의 형태변경이 자주 필요하다면  
꼭 Graphql을 사용해보려고 합니다. 

## 2.2. Repository

### 2.2.1. 업데이트 로직 

[업데이트 로직 예시 ](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L112)

로직은 아래와 같습니다.

1. 요청에 값이 존재하지 않으면 NULL 혹은 "", 0, false로 대체한다.
2. 일대다 관계에서는
   유저가 소유한 ID값과 요청값 ID를 비교하여
	 각 상황에 맞게 생성, 업데이트, 삭제를 진행한다.
3. 다대다 관계에서는
   요청값 id로 모두 대체한다.



### 2.2.2. Patch 로직

[Patch로직 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L292)

로직은 아래와 같습니다.

1. 요청정보가 NULL이 아닌 경우 해당 Field를 업데이트를 한다.
2. 일대다 관계일 경우 요청 id가 존재하지 않으면 생성 존재하면 업데이트를 진행한다.


[업데이트 가능한 컬럼 추출하는 코드](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/utils/repository.go#L52)

요청받은 key가 데이터베이스에서 업데이트할 수 있는지 여부를 확인한 뒤,   
key(컬럼이름)-value(요청값)을 배열에 담아 리턴합니다.

### 2.2.3. List 로직

[List로직 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/recruitment.go#L291)

PostgreSQL에는 jsonb타입이 존재합니다.

json에 있는 key를 불러오기 위해서는 
`column ->> 'key'` 방식으로 코드를 작성해야 합니다.

ent ORM에는 이러한 방식의 쿼리문을 작성할 수 있는 인터페이스가 존재하지 않습니다.

제 프로젝트에서는  
jsonb타입을 많이 사용하지 않고 있기 때문에,  
다른 라이브러리를 찾아보기 보다는  
entORM 인터페이스 내에서 해결하기로 했습니다.  

`테이블 이름, Field이름을 변경해도 로직에 영향을 주지 않기` 중점으로  
프로그램을 작성했습니다.



### 2.2.4. 테스트케이스


[Repository_test](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/user_test.go)

테스트 과정

1. 테스트 시작 시 로컬에서 도커 컨테이너에 DB를 생성
2. 도커 컨테이너에 연결

반복..

3. 각기 다른 테스트 실행에 필요한 데이터 삽입
4. 각각의 테스트 실행
5. 테스트한 데이터 모두 삭제

... 

6. 도커 컨테이너 삭제 

`BeforeTest` 함수에 
실행에 필요한 데이터를 삽입하는 과정을 담았습니다.  

때문에 테스트케이스에서는 테스트할 함수들만 남아있어서  
이해하기 쉬워졌습니다.  

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

DTO, DAO를 분리했습니다.
코드가 길어진다는 단점이 존재하지만,

요청값과 응답값이 명확해진다는 장점과  
요청 응답값 변경이 용이하다는 장점을 포기할 수 없었습니다.

### 2.4.3. 에러처리

[에러 서버코드](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/common/errors.go#L50)

친절한 API를 만들고 싶었습니다.
에러메세지와 에러코드를 상황별로 제공하여
프론트엔드 개발자가 예외처리를 조금 더 편하게 할 수 있도록 하였습니다.

details의 key는 json key값 value에는 어떻게 수정해야 될지 방향을 알려줍니다.

```json
// 응답 예시
HTTP 400 BadReqeust
{
    "code": 2000,
    "message": "유효하지 않은 요청값들이 존재합니다.",
    "details": [
        {
            "email": "email"
        },
        {
            "password": "required"
        },
        {
            "nickname": "required"
        },
        {
            "termAgree": "required"
        }
    ]
}

```

```go
// 응답 예시
HTTP 400 BadReqeust
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}
```

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


