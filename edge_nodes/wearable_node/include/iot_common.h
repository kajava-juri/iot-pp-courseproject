// Common IoT configuration defaults and helpers
#ifndef IOT_COMMON_H
#define IOT_COMMON_H

#ifndef WIFI_NAME
#define WIFI_NAME "TalTech"
#endif

#ifndef WIFI_PASSWORD
#define WIFI_PASSWORD ""
#endif

// Stringize helper for macros coming from -D CPP defines
#ifndef STR_HELPER
#define STR_HELPER(x) #x
#define STR(x) STR_HELPER(x)
#endif

#endif // IOT_COMMON_H
