# Bootstrap for MAC
bootstrap:
	brew install k6

load-test-leaderboard:
	k6 run --out web-dashboard load_test_leaderboard.js

load-test-score:
	k6 run --out web-dashboard load_test.js

.PHONY: bootstrap leaderboard score