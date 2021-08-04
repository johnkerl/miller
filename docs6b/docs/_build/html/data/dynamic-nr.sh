mlr --opprint --from data/small put '
  begin{ @nr1 = 0 }
  @nr1 += 1;
  $nr1 = @nr1
' \
then filter '$x>0.5' \
then put '
  begin{ @nr2 = 0 }
  @nr2 += 1;
  $nr2 = @nr2
'
