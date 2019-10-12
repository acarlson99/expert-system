# Rules
B => A
D + E => B
G + H => F
I + J => G
G => H
L + M => K
O + P => L + N
N => M

# Facts and Queries
=DEIJOP
?AFKP  # AFKP = true
=DEIJP
?AFKP  # AFP = true, K = false

