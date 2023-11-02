all: build

build:
	for i in diagrams/*.dot; do dot -Tpng -Gdpi=300 $$i -o $${i%.dot}.png; done