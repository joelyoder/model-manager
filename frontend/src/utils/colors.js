/**
 * Utility for assigning consistent colors to model tags/badges.
 */

const COLOR_CLASSES = [
    'bg-tag-emerald',
    'bg-tag-indigo',
    'bg-tag-amber',
    'bg-tag-rose',
    'bg-tag-cyan',
    'bg-tag-violet',
    'bg-tag-pink',
    'bg-tag-orange',
    'bg-tag-teal',
    'bg-tag-lime',
    'bg-tag-sky',
    'bg-tag-fuchsia',
];

const FALLBACK_CLASS = 'bg-tag-slate';

// Explicit mappings for known types to ensure they always get the best fitting color
const KNOWN_MAPPINGS = {
    // Model Types
    'checkpoint': 'bg-tag-emerald',
    'lora': 'bg-tag-indigo',
    'textualinversion': 'bg-tag-cyan',
    'embedding': 'bg-tag-cyan',
    'hypernetwork': 'bg-tag-violet',
    'controlnet': 'bg-tag-orange',
    'locon': 'bg-tag-pink',
    'dora': 'bg-tag-rose',
    'poses': 'bg-tag-sky',
    'wildcards': 'bg-tag-amber',
    'vae': 'bg-tag-violet',

    // Base Models
    'sd 1': 'bg-tag-sky',
    'sd 2': 'bg-tag-teal',
    'sdxl': 'bg-tag-rose',
    'sdxl turbo': 'bg-tag-lime',
    'pony': 'bg-tag-pink',
    'flux': 'bg-tag-cyan',
    'illustrious': 'bg-tag-amber',
    'qwen': 'bg-tag-fuchsia',
    'zimage': 'bg-tag-emerald',
};

/**
 * Returns a simple hash number from a string.
 */
function hashCode(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // Convert to 32bit integer
    }
    return Math.abs(hash);
}

/**
 * Returns a CSS class for a given badge text.
 * Uses explicit mappings for known types, otherwise hashes the text 
 * to pick a consistent color from uniform distribution.
 * 
 * @param {string} text - The label text (e.g. "LoRA", "SDXL")
 * @returns {string} - The CSS class string
 */
export function getBadgeColor(text) {
    if (!text) return FALLBACK_CLASS;

    const lower = text.toLowerCase().trim();

    // check for exact or partial matches in known mappings
    // We check keys to see if the text *contains* the key for looser matching
    // (e.g. "Standard LoRA" matching "lora")
    const knownKey = Object.keys(KNOWN_MAPPINGS).find(k => lower.includes(k));
    if (knownKey) {
        return KNOWN_MAPPINGS[knownKey];
    }

    // Fallback: Hash the string to pick a color
    const index = hashCode(lower) % COLOR_CLASSES.length;
    return COLOR_CLASSES[index];
}
