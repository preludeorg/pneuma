#ifdef __WIN32
#include <windows.h>

void RunAgent();

BOOL WINAPI DllMain(
    HINSTANCE _hinstDLL,
    DWORD _fdwReason,
    LPVOID _lpReserved
);

#endif