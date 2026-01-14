---
name: frontend-design
description: Frontend design guidelines for this project. Use when creating or modifying HTML templates, CSS, or any user-facing UI work.
---

# Frontend Design

## How to Use This Skill

1. **Read existing templates first** — `templates/*.html` and `static/css/input.css` are the source of truth for established patterns. Match what's already there.

2. **Use the philosophy below** for new decisions, edge cases, or when you need to understand *why* something is done a certain way.

---

## Design Philosophy

Design should enable the user, not decorate the screen. A "beautiful" interface that hides information, delays action, or prioritizes aesthetics over utility has failed at its primary job.

### Principles

**Information density over whitespace worship.** When someone searches for a list of songs, podcasts, or products—give them a table. Not a grid of giant thumbnails where the actual title gets truncated. Density is not a flaw to be designed away; it's a feature when the user's goal is to scan, compare, and act. Old iTunes understood this. Diskprices.com understands this. Financial terminals understand this.

**Sharpness and contrast where it matters.** Modern UI trends have systematically reduced contrast in the name of "minimalism"—lighter borders, grayed-out text, subtle separators that disappear on certain monitors. This actively harms usability. Edges should be edges. Interactive elements should look interactive. The visual hierarchy should be obvious, not whispered.

**Visible affordances, not animated reveals.** If there's a button to dismiss something, show the button. Don't hide it behind a hover state, then an animation delay, then another animation. Every interaction that requires the user to discover hidden mechanics is friction that accumulates across thousands of interactions.

**Skepticism toward animations.** Animation should serve a functional purpose—communicating state change, orienting the user in a transition. When animation exists for "delight" or "polish," it usually means someone optimized for the demo video rather than the daily user who will see that transition 10,000 times.

**Deliberate rejection of fashion cycles.** Trends come and go: parallax scrolling, skeuomorphism, flat design, neumorphism, glassmorphism. The approach here is consciously unfashionable—drawing instead from patterns that proved themselves over decades. Engineering documentation, scientific instruments, Unix aesthetics. These weren't designed to be trendy; they were designed to work.

### Historical Context

This philosophy isn't nostalgia—it's recognition that certain design patterns were abandoned for reasons other than functionality. The shift from compact list views to giant card grids wasn't because cards are more usable; it's because they fill screens more impressively in screenshots and accommodate less skilled designers who struggle with density. The shift toward hidden UI elements wasn't because users wanted them; it's because "clean" interfaces photograph better.

Look at what professionals use when efficiency matters: Bloomberg terminals, Foobar2000, production DAWs, IDE interfaces (before they started chasing mainstream aesthetics). Look at what has survived without redesign: Craigslist, Hacker News, most effective e-commerce checkout flows.

### The Tension

There's a legitimate counter-argument: consumer software serves diverse audiences, including people who are intimidated by dense interfaces. Fair. But the industry has swung so far toward optimizing for the first-time screenshot that it actively punishes power users and repeat visitors—the people who actually use the product.

The interfaces here serve the person who will use them daily for years, not the person evaluating them for 30 seconds in an app store.
