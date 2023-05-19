docker: env
	@docker build -t gioco-db-migration .
env: dev-env prod-env
dev-env:
ifeq (,$(wildcard etc/.env))
	@cp etc/.example.env etc/.env && echo copy etc/.env done.
endif
prod-env:
ifeq (,$(wildcard etc/.prod.env))
	@cp etc/.example.env etc/.prod.env && echo copy etc/.prod.env done.
endif
.PHONY: docker env dev-env prod-env