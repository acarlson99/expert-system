# Rules
B => A
C => A

# Facts and Queries
=
?A  # A = false

=B
?A  # A = true

=C
?A  # A = true

=BC
?A  # A = true
