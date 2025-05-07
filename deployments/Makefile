HELM = helm

lint:
	$(HELM) lint vapusdata-platform 

install:
	$(HELM) install $(NAME) vapusdata-platform -n $(NAMESPACE) --create-namespace -f $(SETUP_FILE) --set app.dev=$(ISDEV) --set app.namespace=$(NAMESPACE)

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
	$(HELM) dependency update vapusdata-platform
	$(HELM) package vapusdata-platform

release: build
	$(HELM) push vapusdata-platform-helmchart-0.0.2.tgz oci://asia-south1-docker.pkg.dev/vapusdata-beta/vapusdata-ecosystem