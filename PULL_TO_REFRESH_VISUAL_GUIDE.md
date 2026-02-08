# 🎨 Pull-to-Refresh Visual Guide

## 📱 Container Scroll Architecture

```
┌─────────────────────────────────────┐
│  iPhone Screen (h-screen)           │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ Sticky Header (flex-shrink-0)  │ │
│  │ • Safe area top padding        │ │
│  │ • Title, buttons, tabs         │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ Scrollable Container           │ │ ← THIS SCROLLS!
│  │ (flex-1 overflow-y-auto)       │ │
│  │                                │ │
│  │  ┌──────────────────────────┐ │ │
│  │  │ scrollTop = 0            │ │ │ ← TOP
│  │  │ ▼ Pull here = REFRESH ✅ │ │ │
│  │  ├──────────────────────────┤ │ │
│  │  │                          │ │ │
│  │  │ Content...               │ │ │
│  │  │                          │ │ │
│  │  ├──────────────────────────┤ │ │
│  │  │ scrollTop > 0            │ │ │ ← MIDDLE
│  │  │ ▼ Pull here = SCROLL ✅  │ │ │
│  │  ├──────────────────────────┤ │ │
│  │  │                          │ │ │
│  │  │ More content...          │ │ │
│  │  │                          │ │ │
│  │  ├──────────────────────────┤ │ │
│  │  │ pb-24 (96px clearance)   │ │ │
│  │  └──────────────────────────┘ │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ BottomNav (fixed, 50px)        │ │
│  └────────────────────────────────┘ │
│                                      │
│  Home Indicator (34px, system UI)   │
└─────────────────────────────────────┘
```

## 🔄 Pull-to-Refresh States

### State 1: At Top (scrollTop = 0)

```
User at top of scroll container
         │
         ▼
    ┌─────────┐
    │ Content │ ← scrollTop = 0
    │    ▲    │
    │    │    │
    │  Pull   │
    │  Down   │
    └─────────┘
         │
         ▼
   ┌───────────┐
   │ ⬇️ Kéo    │ ← Indicator shows
   │ xuống để  │
   │ làm mới   │
   └───────────┘
         │
         ▼
    Pull > 80px
         │
         ▼
   ┌───────────┐
   │ 🎯 Thả để │ ← Ready to refresh
   │ làm mới   │
   └───────────┘
         │
         ▼
     Release
         │
         ▼
   ┌───────────┐
   │ 🔄 Đang   │ ← Refreshing
   │ tải...    │
   └───────────┘
         │
         ▼
    Data reloads ✅
```

### State 2: Mid-Scroll (scrollTop > 0)

```
User scrolling in middle
         │
         ▼
    ┌─────────┐
    │ Content │ ← scrollTop > 0
    │    ▲    │
    │    │    │
    │  Pull   │
    │  Down   │
    └─────────┘
         │
         ▼
  isPulling = false
         │
         ▼
    ┌─────────┐
    │ Content │ ← Continues scrolling
    │    ▼    │
    │ Scroll  │
    │  Down   │
    └─────────┘
         │
         ▼
   Normal scroll ✅
   NO refresh ✅
```

### State 3: Pull Then Scroll

```
User at top, starts pulling
         │
         ▼
    ┌─────────┐
    │ Content │ ← scrollTop = 0
    │    ▲    │
    │  Pull   │
    └─────────┘
         │
         ▼
   ┌───────────┐
   │ ⬇️ Kéo    │ ← Indicator shows
   │ xuống...  │
   └───────────┘
         │
         ▼
  User scrolls down
         │
         ▼
  scrollTop > 0
         │
         ▼
  Reset pull state
         │
         ▼
    ┌─────────┐
    │ Content │ ← Indicator hides
    │    ▼    │
    │ Scroll  │
    └─────────┘
         │
         ▼
   Normal scroll ✅
   NO refresh ✅
```

## 🔍 Scroll Position Detection

### Old Method (❌ Wrong):

```
┌─────────────────────────────────────┐
│  Page (min-h-screen)                │
│                                      │
│  window.pageYOffset = 0 ← Always 0! │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ Container (overflow-y-auto)    │ │
│  │                                │ │
│  │  scrollTop = 0    ← Top        │ │
│  │  scrollTop = 100  ← Middle     │ │
│  │  scrollTop = 500  ← Bottom     │ │
│  │                                │ │
│  └────────────────────────────────┘ │
│                                      │
└─────────────────────────────────────┘

Problem: Checking window.pageYOffset
         Always returns 0
         Can't detect container scroll position
```

### New Method (✅ Correct):

```
┌─────────────────────────────────────┐
│  Page                                │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ Container (overflow-y-auto)    │ │
│  │                                │ │
│  │  container.scrollTop = 0       │ │ ← Check this!
│  │  ↓                             │ │
│  │  container.scrollTop = 100     │ │ ← Check this!
│  │  ↓                             │ │
│  │  container.scrollTop = 500     │ │ ← Check this!
│  │                                │ │
│  └────────────────────────────────┘ │
│                                      │
└─────────────────────────────────────┘

Solution: Check container.scrollTop
          Accurately detects scroll position
          Works correctly with container scroll
```

## 🎯 Logic Flow

### Touch Start:

```
Touch Start Event
       │
       ▼
Find scroll container
       │
       ├─ Walk up DOM tree
       ├─ Find overflow-y: auto/scroll
       └─ Cache container
       │
       ▼
Get scroll position
       │
       ├─ container.scrollTop (if container)
       └─ window.pageYOffset (if page)
       │
       ▼
   scrollTop = 0?
       │
   ┌───┴───┐
   │       │
  Yes     No
   │       │
   ▼       ▼
Enable  Ignore
Pull    Pull
```

### Touch Move:

```
Touch Move Event
       │
       ▼
  isPulling?
       │
   ┌───┴───┐
   │       │
  Yes     No
   │       │
   ▼       └─> Return
Get scroll position
       │
       ▼
Calculate distance
       │
       ▼
distance > 0 AND scrollTop = 0?
       │
   ┌───┴───┐
   │       │
  Yes     No
   │       │
   ▼       ▼
Show    Reset
Pull    Pull
Indicator
   │       │
   ▼       ▼
Prevent  Allow
Default  Scroll
```

### Touch End:

```
Touch End Event
       │
       ▼
  isPulling?
       │
   ┌───┴───┐
   │       │
  Yes     No
   │       │
   ▼       └─> Return
pullDistance >= threshold?
       │
   ┌───┴───┐
   │       │
  Yes     No
   │       │
   ▼       ▼
Trigger  Reset
Refresh  Pull
   │       │
   ▼       ▼
Show    Hide
Loading Indicator
```

## 📊 Comparison

### Before Fix:

| Scenario | Expected | Actual | Result |
|----------|----------|--------|--------|
| Pull at top | Refresh | Refresh | ✅ OK |
| Pull mid-scroll | Scroll | Refresh | ❌ WRONG |
| Rapid scroll | Scroll | Refresh | ❌ WRONG |

### After Fix:

| Scenario | Expected | Actual | Result |
|----------|----------|--------|--------|
| Pull at top | Refresh | Refresh | ✅ OK |
| Pull mid-scroll | Scroll | Scroll | ✅ OK |
| Rapid scroll | Scroll | Scroll | ✅ OK |

## 🎨 Visual States

### Normal Scroll (scrollTop > 0):

```
┌─────────────────┐
│ Header          │
├─────────────────┤
│ Content         │ ← User sees this
│ ...             │
│ ...             │
│ ▼ Scrolling     │
├─────────────────┤
│ BottomNav       │
└─────────────────┘

No indicator
Normal scroll
```

### Pull at Top (scrollTop = 0):

```
┌─────────────────┐
│ ⬇️ Kéo xuống    │ ← Indicator appears
├─────────────────┤
│ Header          │
├─────────────────┤
│ Content         │ ← Pulled down
│ ...             │
│ ▲ Pulling       │
├─────────────────┤
│ BottomNav       │
└─────────────────┘

Indicator shows
Content pulled
```

### Refreshing:

```
┌─────────────────┐
│ 🔄 Đang tải...  │ ← Loading indicator
├─────────────────┤
│ Header          │
├─────────────────┤
│ Content         │ ← Reloading
│ ...             │
│ ...             │
├─────────────────┤
│ BottomNav       │
└─────────────────┘

Loading animation
Data refreshing
```

## 💡 Key Concepts

### Container Scroll:

```
Page doesn't scroll
    ↓
Container scrolls
    ↓
Check container.scrollTop
    ↓
Not window.pageYOffset
```

### Dynamic Detection:

```
Touch target
    ↓
Walk up DOM
    ↓
Find scrollable parent
    ↓
Cache for reuse
```

### Smart Reset:

```
Start pulling at top
    ↓
User scrolls down
    ↓
scrollTop > 0
    ↓
Reset pull state
    ↓
Allow normal scroll
```

## 🧪 Testing Visual

### Test 1: Pull at Top

```
Before:                After:
┌─────────┐           ┌─────────┐
│ Content │           │ ⬇️ Pull │ ← Indicator
│    ▲    │    →      │ Content │
│  Pull   │           │    ▲    │
└─────────┘           └─────────┘
scrollTop=0           Pulling...

Expected: ✅ Shows indicator, triggers refresh
```

### Test 2: Pull Mid-Scroll

```
Before:                After:
┌─────────┐           ┌─────────┐
│ Content │           │ Content │
│    ▲    │    →      │    ▼    │ ← Scrolls
│  Pull   │           │ Scroll  │
└─────────┘           └─────────┘
scrollTop>0           Scrolling...

Expected: ✅ No indicator, continues scrolling
```

### Test 3: Pull Then Scroll

```
Start:                Middle:               End:
┌─────────┐          ┌─────────┐          ┌─────────┐
│ ⬇️ Pull │    →     │ Content │    →     │ Content │
│ Content │          │    ▼    │          │    ▼    │
│    ▲    │          │ Scroll  │          │ Scroll  │
└─────────┘          └─────────┘          └─────────┘
Pulling...           Scrolling...         Normal

Expected: ✅ Indicator disappears, scrolls normally
```

---

**Visual Guide Created:** February 6, 2026  
**Purpose:** Help understand pull-to-refresh container scroll behavior  
**Audience:** Developers and testers
