HELM = helm

lint:
	$(HELM) lint vapusai

install:
	$(HELM) install $(NAME) vapusai -n $(NAMESPACE) --create-namespace -f $(SETUP_FILE) --set app.dev=$(ISDEV) --set app.namespace=$(NAMESPACE)

uninstall:
	$(HELM) uninstall vapusdata -n $(NAMESPACE)

clean:
	$(HELM) uninstall vapusdata -n $(NAMESPACE)
	kubectl delete clusterrolebinding vapusdata-platform-master-crb
	kubectl delete clusterrole vapusdata-platform-master-cr
	kubectl delete namespace $(NAMESPACE)

upgrade:
	$(HELM) upgrade $(NAME) vapusdata-platform -n $(NAMESPACE)  -f $(SETUP_FILE)

build: lint
	$(HELM) dependency update vapusai
	$(HELM) package vapusai