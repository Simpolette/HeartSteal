# Backend

This is the GO backend built using Clean Architecture, using the template:

```
https://github.com/amitshekhariitbhu/go-backend-clean-architecture
```

To run test:
```
go test -v .\internal\usecase\test\
```

Domain for mocks can be generated using mockery by running mockery.
If not installed, run:
```
go install github.com/vektra/mockery/v2@latest
```

Then:
```
mockery
```