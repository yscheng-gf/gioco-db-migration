docker: env
	@docker build -t gioco-db-migration .

env: dev-env staging-env prod-env

dev-env:
ifeq (,$(wildcard etc/.env))
	@cp etc/.example.env etc/.env && echo copy etc/.env done.
endif

staging-env:
ifeq (,$(wildcard etc/.staging.env))
	@cp etc/.example.env etc/.staging.env && echo copy etc/.staging.env done.
endif

prod-env:
ifeq (,$(wildcard etc/.prod.env))
	@cp etc/.example.env etc/.prod.env && echo copy etc/.prod.env done.
endif
.PHONY: docker env dev-env prod-env