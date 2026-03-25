import React, { useState, useRef, useCallback } from 'react';
import HeroSection from '../components/HeroSection';
import GeneratorForm from '../components/GeneratorForm';

const HomePage: React.FC = () => {
    const [formInteracted, setFormInteracted] = useState(false);
    const generatorRef = useRef<HTMLDivElement>(null);

    const handleInteract = useCallback(() => {
        setFormInteracted(true);
    }, []);

    const handleOpenGenerator = useCallback(() => {
        generatorRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, []);

    return (
        <div className="main-centered">
            <HeroSection collapsed={formInteracted} onOpenGenerator={handleOpenGenerator} />
            <div ref={generatorRef}>
                <GeneratorForm onInteract={handleInteract} />
            </div>
        </div>
    );
};

export default HomePage;
