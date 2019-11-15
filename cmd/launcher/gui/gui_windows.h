#include <windows.h>
#include <Commctrl.h>

int centerWindow(ULONG_PTR ulpWindowHandle) {
    HANDLE windowHandle = (HANDLE)ulpWindowHandle;
    HMONITOR monitor = MonitorFromWindow(windowHandle, MONITOR_DEFAULTTOPRIMARY);
    if (monitor == NULL) {
        return 1;
    }
    MONITORINFO monitorInfo = { .cbSize = sizeof(MONITORINFO) };
    if (GetMonitorInfo(monitor, &monitorInfo) == 0) {
        return 2;
    }
    RECT rect;
    if (GetWindowRect(windowHandle, &rect) == 0) {
        return 3;
    }
    LONG windowWidth = rect.right - rect.left;
    LONG windowHeight = rect.bottom - rect.top;
    LONG newX = (monitorInfo.rcWork.right - windowWidth) / 2;
    LONG newY = (monitorInfo.rcWork.bottom - windowHeight) / 2;
    SetWindowPos(windowHandle, HWND_TOP, newX, newY, 0, 0, SWP_NOACTIVATE | SWP_NOOWNERZORDER | SWP_NOSIZE);
}

HICON largeIcon = NULL;
HICON smallIcon = NULL;

UINT loadIcons(LPCWSTR binaryPath) {
    HICON* largeIconPtr = calloc(1, sizeof(HICON*));
    HICON* smallIconPtr = calloc(1, sizeof(HICON*));
    UINT extractedIconCount = ExtractIconExW(binaryPath, 0, largeIconPtr, smallIconPtr, 1);
    largeIcon = *largeIconPtr;
    smallIcon = *smallIconPtr;
    free(largeIconPtr);
    free(smallIconPtr);
    return extractedIconCount;
}

void applyIconToWindow(ULONG_PTR ulpWindowHandle) {
    HANDLE windowHandle = (HANDLE)ulpWindowHandle;
    if (largeIcon != NULL) {
        SendMessage((HWND)windowHandle, WM_SETICON, ICON_BIG, (LPARAM)largeIcon);
    }
    if (smallIcon != NULL) {
        SendMessage((HWND)windowHandle, WM_SETICON, ICON_SMALL, (LPARAM)smallIcon);
    }
}

void applyWindowStyle(ULONG_PTR ulpWindowHandle) {
    HANDLE windowHandle = (HANDLE)ulpWindowHandle;
    SetWindowLong(windowHandle, GWL_STYLE, GetWindowLong(windowHandle, GWL_STYLE)&~(WS_SIZEBOX|WS_MAXIMIZEBOX));
}

void setProgressBarState(ULONG_PTR ulpBarHandle, int progressState) {
    HWND barHWND = (HWND)ulpBarHandle;
    SendMessage(barHWND, 1040, (WPARAM)progressState, 0);
}