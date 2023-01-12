# onthemat_backend


데모 (미완성)

https://user-images.githubusercontent.com/97140962/210781527-ece2bf08-1c96-47fa-b04b-9a58e513239e.mov

주요 기술적 과제 

- 확장 가능한 서버 아키텍처
- 엘라스틱 서치를 활용하여 검색시스템 구축
- RFC규격에 맞는 POST, PUT, DELETE, PATCH 구축
- Mock 코드 자동생성, 유닛테스트 구축
- jwt 토큰 인증, oauth
- fiber, entGo (장점 : 코드제너레이팅을 통해, 컴파일 단계에서 에러를 잡을 수 있음) 학습 후 적용
- 친절한 에러 메세지



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
    <img src="https://img.shields.io/badge/Elastic-green?style=flat&logo=ElasticSearch&logoColor=005571"/>
</div>
<br>

fiber를 선택한 이유 : 성능

![image](https://user-images.githubusercontent.com/97140962/205558298-df3012cd-5f72-43a6-a158-1d987105198c.png)

entGO를 선택한 이유 : 컴파일 단계에서 에러를 잡을 수 있음. 

postgres를 선택한 이유 : fiber는 fasthttp기반으로 설계가 되어 있음, fasthttp와 postgresql의 조합이 빠른 속도를 냄.


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

URL : http://13.125.48.238:3000/

<img width="1680" alt="스크린샷 2022-12-04 오후 3 54 02" src="https://user-images.githubusercontent.com/97140962/205478479-f0c7fd16-e8fd-4590-81f8-a1b8c386cfa6.png">


# 2. 기술 상세

## 2.1. REST API


GET, POST, PUT, PACTH, DELETE 5가지를 사용합니다.

RFC문서의 내용에 맞게 설계했습니다.

> The PUT method requests that the state of the target resource be created or 
> replaced with the state defined by the representation enclosed in the request message payload.

> This specification defines the new HTTP/1.1 [RFC2616] method, 
> PATCH, which is used to apply partial modifications to a resource.



## 2.2. Repository



[Put 로직 예시 ](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L112)

[Patch로직 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/teacher.go#L292)

[업데이트 가능한 컬럼 추출하는 코드](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/utils/repository.go#L52)

[List로직 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/recruitment.go#L291)


### 2.2.4. 테스트케이스


[Repository_test](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/repository/user_test.go)

테스트 과정

1. 테스트 시작 시 로컬에서 도커 컨테이너에 DB를 생성
2. 도커 컨테이너에 연결

...반복...

3. 테스트 실행에 필요한 데이터 삽입
4. 테스트 실행
5. 테스트한 데이터 모두 삭제

... 


6. 도커 컨테이너 삭제


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

Usecase에서는 repostiroy, service 모듈이 사용됩니다.
repository, service 모듈들은 유닛테스트가 완료됐습니다.

mock을 사용하여 usecase에서만 작성된 로직만 테스트할 수 있도록 하였습니다.

## 2.4. handler

[코드 예시](https://github.com/beardfriend/onthemat_backend/blob/main/internal/app/delivery/http/academy_handler.go#L30)

응답받은 요청을 파싱, 검증하는 과정을 담당합니다.


