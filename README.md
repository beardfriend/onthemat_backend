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
│   └── app             // 메인 어플리케이션
├── configs             // 어플리케이션 설정파일
│
├── internal
│    └── app 
│         ├── delivery            
│         │           ├── http
│         │           └── middleware
│         ├── config                  
│         │
│         ├── common
│         │
│         ├── infrastructure
│         │
│         ├── model
│         │
│         ├── repository
│         │
│         ├── service
│         │
│         ├── usecase
│         │
│         ├── transport
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
│    ├── google
│    │
│    ├── kakao
│    │
│    ├── naver
│    │
│    ├── openapi
│    │
│    ├── validatorx
```

## 1.3. 아키텍처


### 1.3.1. 서버 아키텍처

의존성을 최대한 낮추기 위해 노력했습니다.

![서버 아키텍처](https://user-images.githubusercontent.com/97140962/202600851-884abaad-c12c-4f7e-8b23-715dee475e5c.jpg)

## 1.4. API 명세 
![api문서](https://user-images.githubusercontent.com/97140962/201019708-08588b56-8304-4a77-946a-cf67e443a7a5.png)



