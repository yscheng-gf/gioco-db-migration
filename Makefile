docker:
	@docker build -t gioco-db-migration .
env:
ifeq (,$(wildcard etc/.env))
	@cp etc/.env.example etc/.env && echo copy etc/.env done.
endif
ifeq (,$(wildcard etc/.prod.env))
	@cp etc/.env.example etc/.prod.env && echo copy etc/.prod.env done.
endif