# **SMS-Health-Check**

> ### 대덕소프트웨어마이스터고등학교 [**SMS(학교 지원 시스템)**](https://github.com/DMS-SMS) 서버 **Health Check 서비스**

<br>

---

## **INDEX**
### [**1. SMS Health Check란?**](#SMS-Health-Check란?)
### [**2. Health Check 종류**](#Health-Check-종류)
### [**3. 패키지 의존 흐름 및 그래프**](#패키지-의존성-흐름-및-그래프)
### [**4. 프로젝트 구조**](#프로젝트-구조)
### [**5. 부가 기능**](#부가-기능)

<br>

---

## **SMS Health Check란?**
- **SMS**(School Management System)는 전공 동아리 [**DMS**](https://github.com/DSM-DMS)에서 개발하여 현재 운영하고 있는 **학교 지원 시스템**입니다.
- **SMS Health Check**는 SMS 서버를 **운영중인 환경**과 서버를 **구성중인 여러 서비스**들의 상태를 **주기적으로 관리**하는 플러그인 식의 서비스입니다.
- 특정 서비스의 **상태 확인 결과**가 특정 기준을 통해 **정상적이지 않다고 판단**이 되면 해당 서비스의 **상태 회복을 위한 작업을 수행**합니다.

<br>

![godepgraph1](https://user-images.githubusercontent.com/48676834/113800510-08960180-9792-11eb-8c5d-a5650ab0799b.png)

![godepgraph2](https://user-images.githubusercontent.com/48676834/113800517-0df34c00-9792-11eb-90ab-d048f3b847a1.png)
