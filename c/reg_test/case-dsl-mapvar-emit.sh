# ----------------------------------------------------------------
announce MAPVAR EMIT

run_mlr -n put 'end {
  @a[2][3] = 4;
  b[2][3] = 8;
  emit  (@a, b, {2:{3:12}});
  emitp (@a, b, {2:{3:12}});
  emit {};
  emit  (@a, b, {2:{3:12}}), "t";
  emitp (@a, b, {2:{3:12}}), "t";
  emit {};
  emit  (@a, b, {2:{3:12}}), "t", "u";
  emitp (@a, b, {2:{3:12}}), "t", "u";
}'

run_mlr -n put 'end {
  @a[2][3] = 4;
  b[2][3] = 8;
  emitp (@a, b, {2:{3:12}});
  emit  (@a, b, {2:{3:12}});
  emit {};
  emitp (@a, b, {2:{3:12}}), "t";
  emit  (@a, b, {2:{3:12}}), "t";
  emit {};
  emitp (@a, b, {2:{3:12}}), "t", "u";
  emit  (@a, b, {2:{3:12}}), "t", "u";
}'
