# Visual Design System for Jambe Verte

## Overview
A clean, minimalist blog design inspired by Mitchell Hashimoto's blog and Anthropic's readability, with a sage green color palette and approachable aesthetic.

## Color Palette

### Primary Colors
- **Background**: `#FDFCF8` - Warm off-white (slightly warmer than pure white)
- **Text**: `#1F2937` - Dark charcoal gray
- **Primary Green**: `#7C9885` - Sage green (main brand color)
- **Darker Green**: `#5A6B61` - Forest sage (for hover states and emphasis)
- **Accent Pink**: `#E88D9D` - Soft dusty pink (for special highlights)

### Supporting Colors
- **Muted Text**: `#6B7280` - Medium gray (for dates, metadata)
- **Border**: `#E5E7EB` - Light gray
- **Code Background**: `#F3F4F6` - Cool light gray
- **Selection**: `#7C988520` - Sage green at 20% opacity

## Typography

### Font Stack
```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
```

### Type Scale
- **Body**: 16px (1rem), line-height 1.75
- **Small**: 14px (0.875rem)
- **H1**: 36px (2.25rem), font-weight 700
- **H2**: 30px (1.875rem), font-weight 600
- **H3**: 24px (1.5rem), font-weight 600
- **H4**: 20px (1.25rem), font-weight 600

### Blog Post Typography
- **Paragraph spacing**: 1.5rem bottom margin
- **Max line length**: 65-70 characters (approximately 650px)
- **Letter spacing**: -0.01em for headings

## Layout Structure

### Desktop (≥768px)
```
+------------------+------------------------+
|                  |                        |
|    Sidebar       |    Main Content        |
|    (250px)       |    (max 650px)         |
|                  |                        |
|  - Home          |                        |
|  - Posts         |                        |
|  - About         |                        |
|                  |                        |
+------------------+------------------------+
```

### Mobile (<768px)
```
+--------------------------------+
|         Top Navigation         |
|    Home | Posts | About        |
+--------------------------------+
|                                |
|        Main Content            |
|         (padding)              |
|                                |
+--------------------------------+
```

## Component Specifications

### Navigation
- **Desktop Sidebar**:
  - Fixed position
  - Width: 250px
  - Background: transparent
  - Border-right: 1px solid `#E5E7EB`
  - Padding: 2rem
  
- **Navigation Items**:
  - Font-size: 16px
  - Color: `#6B7280` (default)
  - Active/hover: `#1F2937`
  - Active indicator: 3px left border in `#7C9885`
  - Padding: 0.75rem 1rem
  - Transition: all 0.15s ease

- **Mobile Navigation**:
  - Sticky top bar
  - Background: `#FDFCF8` with subtle shadow
  - Height: 60px
  - Horizontal layout with equal spacing

### Content Area
- **Container**:
  - Max-width: 650px
  - Padding: 2rem (desktop), 1.5rem (mobile)
  - Margin: 0 auto

### Blog Post Cards
- **Card Structure**:
  - No border or background
  - Separator: 1px solid `#E5E7EB` between posts
  - Padding: 2rem 0
  
- **Post Title**:
  - Font-size: 24px
  - Color: `#1F2937`
  - Hover: `#5A6B61` (darker green)
  - Transition: color 0.15s ease
  
- **Metadata**:
  - Date: `#6B7280`, font-size 14px
  - Tags: Background `#7C988520`, text `#5A6B61`
  - Tag hover: Background `#7C988540`

### Spacing System
Using Tailwind's default spacing scale:
- xs: 0.5rem (8px)
- sm: 1rem (16px)
- md: 1.5rem (24px)
- lg: 2rem (32px)
- xl: 3rem (48px)

## Interactive States

### Links
- Default: `#7C9885` (sage green)
- Hover: `#5A6B61` (darker green)
- Visited: `#7C9885` (same as default)
- Underline on hover only

### Buttons (if needed)
- Primary: Background `#7C9885`, text white
- Hover: Background `#5A6B61`
- Focus ring: 2px offset, `#7C988550`

### Form Elements (if needed)
- Border: `#E5E7EB`
- Focus border: `#7C9885`
- Background: white
- Focus ring: 2px `#7C988550`

## Tailwind Configuration

### Extend colors:
```javascript
colors: {
  'sage': {
    50: '#F5F7F5',
    100: '#E8ECE9',
    200: '#C8D4CC',
    300: '#A7BBAF',
    400: '#7C9885',  // Primary
    500: '#5A6B61',  // Darker
    600: '#475650',
    700: '#384441',
    800: '#2B3432',
    900: '#1F2523',
  },
  'blush': {
    400: '#E88D9D',  // Accent pink
    500: '#E17589',
  },
  'cream': '#FDFCF8',  // Background
}
```

### Key Tailwind Classes
- Background: `bg-cream`
- Text: `text-gray-800`
- Primary green: `text-sage-400`, `bg-sage-400`
- Darker green: `text-sage-500`, `bg-sage-500`
- Pink accent: `text-blush-400`

## Implementation Notes

1. **Font Loading**: Load Inter from Google Fonts with weights 400, 600, 700
2. **Responsive**: Mobile-first approach using Tailwind's responsive prefixes
3. **Accessibility**: Maintain WCAG AA contrast ratios
4. **Performance**: Minimize custom CSS, leverage Tailwind utilities
5. **Dark Mode**: Not included in initial design (can be added later)

## Example Component Structure

### Desktop Sidebar Item (Active)
```html
<a class="block px-4 py-3 text-gray-800 border-l-4 border-sage-400 bg-sage-50 font-medium">
  Posts
</a>
```

### Blog Post Card
```html
<article class="py-8 border-b border-gray-200 last:border-0">
  <h2 class="mb-2">
    <a class="text-2xl font-semibold text-gray-800 hover:text-sage-500 transition-colors">
      Post Title
    </a>
  </h2>
  <time class="text-sm text-gray-500">January 13, 2025</time>
  <p class="mt-3 text-gray-600 leading-relaxed">
    Post description or excerpt...
  </p>
  <div class="mt-3 flex gap-2">
    <span class="px-3 py-1 text-sm bg-sage-400/20 text-sage-600 rounded-full">
      golang
    </span>
  </div>
</article>
```