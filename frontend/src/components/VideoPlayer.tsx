'use client';

type Props = {
  src: string;
  poster?: string;
  className?: string;
};

export default function VideoPlayer({ src, poster, className }: Props) {
  const isFile = /\.(mp4|webm|m3u8|mov)(\?|$)/i.test(src);
  if (isFile) {
    return (
      <video
        src={src}
        poster={poster}
        controls
        playsInline
        preload="metadata"
        controlsList="nodownload noplaybackrate nopictureinpicture"
        onContextMenu={(e) => e.preventDefault()}
        className={className || 'w-full h-full bg-black'}
      />
    );
  }
  return (
    <iframe
      src={src}
      className={className || 'w-full h-full'}
      allow="autoplay; fullscreen; encrypted-media"
      allowFullScreen
      referrerPolicy="no-referrer"
      sandbox="allow-scripts allow-same-origin allow-presentation"
    />
  );
}
