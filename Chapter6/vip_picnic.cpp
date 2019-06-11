#include <iostream>

using namespace std;

#define N 10

bool assign[N] = { false };
bool areFriends[N][N] = { { false } };

int findPair(int n) {
	int candidate, mate;
	int cases;

	if (n < 2) {
		for (candidate = 0; candidate < N; candidate++) {
			if (!assign[candidate]) return 0;
		}
		return 1;
	}

	for (candidate = 0; candidate < N; candidate++) {
		if (!assign[candidate]) {
			assign[candidate] = true;
			break;
		}
	}

	for (mate = candidate+1; mate < N; mate++) {
		if (areFriends[candidate][mate] && !assign[mate]) {
			assign[mate] = true;
			cases += findPair(n-2);
			assign[mate] = false;
		}
	}

	assign[candidate] = false;

	return cases;
}

int main() {
	int n, m;
	int x, y;

	cin >> n; cin >> m;

	for (int i = n; i < N; i++) {
		assign[i] = true;
	}

	for (int i = 0; i < m; i++) {
		cin >> x; cin >> y;
		if (x < y) {
			areFriends[x][y] = true;
		} else {
			areFriends[y][x] = true;
		}
	}

	cout << "Result: " << findPair(n) << endl;
}
