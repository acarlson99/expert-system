# Rules
A | B + C => E
(F | G) + H => E

# Facts and Queries
=
?E  # E = false

=A
?E  # E = true

=B
?E  # E = false

=C
?E  # E = false

=AC
?E  # E = true

=BC
?E  # E = true

=F
?E  # E = false

=G
?E  # E = false

=H
?E  # E = false

=FH
?E  # E = true

=GH
?E  # E = true
