class MyClass {
  private:
    char* a;
  public:
  MyClass(char* a) {
    this->a = strdup(a);
  }
  ~MyClass() {
    free(a);
  }
  int myMethod(char* b) {
    return strlen(a) + strlen(b);
  }
};
int main(void) {
  MyClass* myObj = new MyClass("hello");
  int x = myObj.myMethod("world");
  return 0;
}
