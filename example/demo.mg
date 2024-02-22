edge(/a, /b).
edge(/b, /c).
edge(/c, /d).

reachable(X, Y) :- edge(X, Y).
reachable(X, Z) :- edge(X, Y), reachable(Y, Z).

