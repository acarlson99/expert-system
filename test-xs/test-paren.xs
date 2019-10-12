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
# The correction sheet says false, but that would not make sense
# because + has a greater precedence than |, so the first expression would be
# (A | B) + C => E
# which would be true as long as C and either A or B are true

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
