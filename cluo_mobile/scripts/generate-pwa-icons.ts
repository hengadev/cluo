import sharp from 'sharp';
import { join } from 'path';

// Shared "C" glyph path used across all icon variants.
function cGlyph(size: number, stroke: string) {
	const sw = size * 0.08;
	return `
		<path d="M ${size * 0.72} ${size * 0.28}
		         Q ${size * 0.72} ${size * 0.22} ${size * 0.60} ${size * 0.22}
		         L ${size * 0.38} ${size * 0.22}
		         Q ${size * 0.22} ${size * 0.22} ${size * 0.22} ${size * 0.38}
		         L ${size * 0.22} ${size * 0.62}
		         Q ${size * 0.22} ${size * 0.78} ${size * 0.38} ${size * 0.78}
		         L ${size * 0.60} ${size * 0.78}
		         Q ${size * 0.72} ${size * 0.78} ${size * 0.72} ${size * 0.72}"
			  fill="none"
			  stroke="${stroke}"
			  stroke-width="${sw}"
			  stroke-linecap="round"
			  stroke-linejoin="round" />`;
}

// Maskable/home-screen icons: rounded rect with transparent corners.
async function generateIcon(size: number, outputPath: string) {
	const svg = `
	<svg width="${size}" height="${size}" xmlns="http://www.w3.org/2000/svg">
		<rect width="${size}" height="${size}" fill="#ffffff" rx="${size * 0.15}" />
		${cGlyph(size, '#1a1a1a')}
	</svg>`;
	await sharp(Buffer.from(svg)).png().toFile(outputPath);
	console.log(`Generated ${outputPath}`);
}

// Staging variant: white C on dark rounded rect.
async function generateStagingIcon(size: number, outputPath: string) {
	const svg = `
	<svg width="${size}" height="${size}" xmlns="http://www.w3.org/2000/svg">
		<rect width="${size}" height="${size}" fill="#1a1a1a" rx="${size * 0.15}" />
		${cGlyph(size, '#ffffff')}
	</svg>`;
	await sharp(Buffer.from(svg)).png().toFile(outputPath);
	console.log(`Generated ${outputPath}`);
}

// Apple touch icon: full opaque square (no transparency, no rounding).
// iOS applies its own mask/corners, so the source must be square + opaque.
async function generateAppleTouchIcon(size: number, outputPath: string, staging = false) {
	const bg = staging ? '#1a1a1a' : '#ffffff';
	const stroke = staging ? '#ffffff' : '#1a1a1a';
	const svg = `
	<svg width="${size}" height="${size}" xmlns="http://www.w3.org/2000/svg">
		<rect width="${size}" height="${size}" fill="${bg}" />
		${cGlyph(size, stroke)}
	</svg>`;
	// Flatten onto an opaque background → drop the alpha channel entirely
	await sharp(Buffer.from(svg)).flatten({ background: bg }).png().toFile(outputPath);
	console.log(`Generated ${outputPath}`);
}

async function main() {
	const staticDir = join(process.cwd(), 'static');

	// Production icons: dark C on white
	await generateIcon(192, join(staticDir, 'icon-192.png'));
	await generateIcon(512, join(staticDir, 'icon-512.png'));

	// Staging icons: white C on dark
	await generateStagingIcon(192, join(staticDir, 'icon-staging-192.png'));
	await generateStagingIcon(512, join(staticDir, 'icon-staging-512.png'));

	// Apple touch icons (opaque 180×180)
	await generateAppleTouchIcon(180, join(staticDir, 'icon-180.png'));
	await generateAppleTouchIcon(180, join(staticDir, 'icon-staging-180.png'), true);

	// Favicons (production)
	await generateIcon(32, join(staticDir, 'favicon-32x32.png'));
	await generateIcon(16, join(staticDir, 'favicon-16x16.png'));

	console.log('PWA icons generated successfully!');
}

main().catch(console.error);
