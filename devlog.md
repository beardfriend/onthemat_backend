# 해시태그 기능

## 기능 구현 단계에서 방법 선택

해시태그 생성 업데이트 삭제 로직의 경우를 따져보자.

1. 유저가 갖고 있는 값이 빈 배열일 경우 
   1. 생성(CreateMany)
  


2. 유저가 갖고 있는 값이 있는 경우
   1. 요청값이 기존 값과 일치하지 않는 경우 
  (
		기존 값에는 존재하지만 요청에 없을 경우 삭제, 
		기존값에는 존재하지 않지만 요청에 있을 경우 생성
	)
   1. 요청값이 []인 경우 (Delete ALL)
   2. 요청값이 기존 값과 같은 경우 (No action)


API를 2개 DeleteMany, CreateMany로 나눈다 [✔️]

API를 1개로 할 경우, 기존 값 모두 삭제 -> Request에 있는 값 모두 생성



## 로직구현 Usecase[✔️] vs Repository[✔️]

현재 구조는
Usecase에서 Repository,Service에서 재료를 가져와서
요리를 하는 방식이다.

트렌젝션 Usecase[] vs Repository[✔️]
리스팅 로직 Usecase[] vs Repository[✔️]
비즈니스 로직 Usecase[✔️] vs Repository[✔️]


## JSON형태로 저장[]  vs 새로운 테이블[✔️] (one to Many)

- JSON으로 저장할 경우, 만약에 아쉬탕가를 등록한 선생님을 찾고 싶으면, 찾기가 참 까다롭다. ( 엘라스틱으로 해결 가능? -> 어쩄든 비용이니까..)
- 테이블을 하나 만들어서 하는 편이 좋아보임.

## Reference 요가를 어떤 식으로 참조할 것인지

# ToDO


[ ] 요가 Raw요가랑 합칠 때 정렬 순서를 사용자가 정할 수 있도록 column 추가

[ ] Recruit CRUD

[ ] 지원하기

[ ] 최근 활동내역 , 알림기능

[ ] CreatedAt UpdatedAt GoType바꾸기